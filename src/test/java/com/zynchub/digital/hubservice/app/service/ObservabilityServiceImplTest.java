package com.zynchub.digital.hubservice.app.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.hamcrest.Matchers.allOf;
import static org.hamcrest.Matchers.containsString;
import static org.hamcrest.Matchers.not;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;

import com.zynchub.digital.hubservice.app.config.ObservabilityClientProperties;
import com.zynchub.digital.hubservice.app.dto.AdminObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.ConnectorHealthListResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityInventoryItemResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.TenantHealthListResponseDto;
import com.zynchub.digital.hubservice.app.service.impl.ObservabilityServiceImpl;
import java.time.Instant;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import org.springframework.web.client.RestClient;

class ObservabilityServiceImplTest {

  @Test
  void getOverview_forwards_time_window_period_and_tenant() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String baseUrl = "https://observability.test";

    server.expect(requestTo(allOf(
            containsString(baseUrl + "/api/v1/observability/overview"),
            containsString("from="),
            containsString("to="),
            containsString("periodSeconds=300"),
            containsString("tenantId=tenant-a"))))
        .andRespond(withSuccess(
            "{\"from\":\"2026-03-09T00:00:00Z\",\"to\":\"2026-03-10T00:00:00Z\",\"totalRequests\":42}",
            MediaType.APPLICATION_JSON));

    ObservabilityService service =
        new ObservabilityServiceImpl(
            builder, new ObservabilityClientProperties.Properties(baseUrl));

    ObservabilityOverviewResponseDto response =
        service.getOverview(
            Instant.parse("2026-03-09T00:00:00Z"),
            Instant.parse("2026-03-10T00:00:00Z"),
            300,
            "tenant-a");

    assertThat(response).isNotNull();
    assertThat(response.getTotalRequests()).isEqualTo(42);
    server.verify();
  }

  @Test
  void getInventory_omits_tenant_query_param_when_not_provided() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String baseUrl = "https://observability.test";

    server.expect(requestTo(allOf(
            containsString(baseUrl + "/api/v1/observability/integrations/inventory"),
            containsString("from="),
            containsString("to="),
            not(containsString("tenantId=")))))
        .andRespond(withSuccess(
            "[{\"flowId\":\"flow-1\",\"apiId\":\"api-1\",\"operationId\":\"op-1\",\"connectorId\":\"conn-1\",\"lastSeenAt\":\"2026-03-10T00:00:00Z\",\"requestCount\":7,\"successRate\":99.0,\"statusLight\":\"green\"}]",
            MediaType.APPLICATION_JSON));

    ObservabilityService service =
        new ObservabilityServiceImpl(
            builder, new ObservabilityClientProperties.Properties(baseUrl));

    List<ObservabilityInventoryItemResponseDto> response =
        service.getInventory(
            Instant.parse("2026-03-09T00:00:00Z"),
            Instant.parse("2026-03-10T00:00:00Z"),
            null);

    assertThat(response).hasSize(1);
    assertThat(response.get(0).getFlowId()).isEqualTo("flow-1");
    server.verify();
  }

  @Test
  void getAdminOverview_forwards_period_tenant_and_connector() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String baseUrl = "https://observability.test";

    server.expect(requestTo(allOf(
            containsString(baseUrl + "/api/v1/admin/observability/overview"),
            containsString("periodSeconds=300"),
            containsString("tenantId=tenant-a"),
            containsString("connectorId=dvla"))))
        .andRespond(withSuccess(
            "{\"from\":\"2026-03-09T00:00:00Z\",\"to\":\"2026-03-10T00:00:00Z\",\"totalRequests\":1000,\"successRate\":99.9}",
            MediaType.APPLICATION_JSON));

    ObservabilityService service =
        new ObservabilityServiceImpl(
            builder, new ObservabilityClientProperties.Properties(baseUrl));

    AdminObservabilityOverviewResponseDto response =
        service.getAdminOverview(
            Instant.parse("2026-03-09T00:00:00Z"),
            Instant.parse("2026-03-10T00:00:00Z"),
            300,
            "tenant-a",
            "dvla");

    assertThat(response).isNotNull();
    assertThat(response.getTotalRequests()).isEqualTo(1000);
    server.verify();
  }

  @Test
  void getAdminTenantHealth_forwards_status_paging_and_maps_response() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String baseUrl = "https://observability.test";

    server.expect(requestTo(allOf(
            containsString(baseUrl + "/api/v1/admin/observability/tenants"),
            containsString("status=critical"),
            containsString("page=2"),
            containsString("size=25"))))
        .andRespond(withSuccess(
            "{\"total\":1,\"items\":[{\"tenantId\":\"tenant-a\",\"totalRequests\":50,\"successRate\":92.0,\"failureRate\":8.0,\"quotaUsedPercent\":77.0,\"throttledCount\":3,\"timeoutCount\":1,\"status\":\"critical\"}]}",
            MediaType.APPLICATION_JSON));

    ObservabilityService service =
        new ObservabilityServiceImpl(
            builder, new ObservabilityClientProperties.Properties(baseUrl));

    TenantHealthListResponseDto response =
        service.getAdminTenantHealth(
            Instant.parse("2026-03-09T00:00:00Z"),
            Instant.parse("2026-03-10T00:00:00Z"),
            "critical",
            2,
            25);

    assertThat(response).isNotNull();
    assertThat(response.getTotal()).isEqualTo(1);
    assertThat(response.getItems()).hasSize(1);
    assertThat(response.getItems().get(0).getTenantId()).isEqualTo("tenant-a");
    server.verify();
  }

  @Test
  void getAdminConnectorHealth_forwards_status_paging_and_maps_response() {
    RestClient.Builder builder = RestClient.builder();
    MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
    String baseUrl = "https://observability.test";

    server.expect(requestTo(allOf(
            containsString(baseUrl + "/api/v1/admin/observability/connectors"),
            containsString("status=degraded"),
            containsString("page=1"),
            containsString("size=10"))))
        .andRespond(withSuccess(
            "{\"total\":1,\"items\":[{\"connectorId\":\"dvla\",\"totalRequests\":999,\"successRate\":95.5,\"failureRate\":4.5,\"p95LatencyMs\":1200.0,\"activeTenants\":14,\"status\":\"degraded\"}]}",
            MediaType.APPLICATION_JSON));

    ObservabilityService service =
        new ObservabilityServiceImpl(
            builder, new ObservabilityClientProperties.Properties(baseUrl));

    ConnectorHealthListResponseDto response =
        service.getAdminConnectorHealth(
            Instant.parse("2026-03-09T00:00:00Z"),
            Instant.parse("2026-03-10T00:00:00Z"),
            "degraded",
            1,
            10);

    assertThat(response).isNotNull();
    assertThat(response.getTotal()).isEqualTo(1);
    assertThat(response.getItems()).hasSize(1);
    assertThat(response.getItems().get(0).getConnectorId()).isEqualTo("dvla");
    server.verify();
  }
}
