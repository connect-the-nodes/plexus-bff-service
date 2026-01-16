package com.zynchub.digital.hubservice.app.feature.flags.mapper;

import static org.assertj.core.api.Assertions.assertThat;

import com.zynchub.digital.hubservice.app.feature.flags.dto.FeatureDto;
import com.zynchub.digital.hubservice.app.feature.flags.model.FeatureFlag;
import java.util.List;
import org.junit.jupiter.api.Test;

class FeatureMapperTest {

  @Test
  void maps_parent_to_self_when_blank() {
    FeatureMapper mapper = new FeatureMapper();
    FeatureDto dto = new FeatureDto("FEATURE_ONE", true, "desc", " ");

    List<FeatureFlag> flags = mapper.mapFromDto(List.of(dto));

    assertThat(flags).hasSize(1);
    FeatureFlag flag = flags.get(0);
    assertThat(flag.getName()).isEqualTo("FEATURE_ONE");
    assertThat(flag.getParent()).isEqualTo("FEATURE_ONE");
    assertThat(flag.getDescription()).isEqualTo("desc");
    assertThat(flag.isEnabled()).isTrue();
  }
}
