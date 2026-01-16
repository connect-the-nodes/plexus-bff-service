package com.zynchub.digital.hubservice.app.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;

import com.zynchub.digital.hubservice.app.exception.FeatureNotFoundException;
import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import com.zynchub.digital.hubservice.app.service.impl.RemoteFeaturesRetrieverImpl;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import org.springframework.web.client.RestClient;

class RemoteFeaturesRetrieverImplTest {

  @Test
  void retrieves_and_maps_features_from_remote_yaml() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String url = "https://features.test/flags";

    String yaml =
        "features:\n"
            + "  - name: REMOTE_PARENT\n"
            + "    enabled: true\n"
            + "    description: Remote parent\n"
            + "    parent: REMOTE_PARENT\n"
            + "  - name: REMOTE_CHILD\n"
            + "    enabled: true\n"
            + "    parent: REMOTE_PARENT\n";

    server.expect(requestTo(url)).andRespond(withSuccess(yaml, MediaType.TEXT_PLAIN));

    FeaturesRetriever retriever =
        new RemoteFeaturesRetrieverImpl(builder.build(), url, new FeatureMapper());

    List<FeatureFlag> flags = retriever.retrieveFeatures();

    assertThat(flags).extracting(FeatureFlag::getName).contains("REMOTE_PARENT", "REMOTE_CHILD");
    server.verify();
  }

  @Test
  void throws_when_remote_yaml_is_invalid() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String url = "https://features.test/flags";

    server.expect(requestTo(url)).andRespond(withSuccess("not: yaml: [", MediaType.TEXT_PLAIN));

    FeaturesRetriever retriever =
        new RemoteFeaturesRetrieverImpl(builder.build(), url, new FeatureMapper());

    assertThatThrownBy(retriever::retrieveFeatures)
        .isInstanceOf(FeatureNotFoundException.class);
    server.verify();
  }
}
