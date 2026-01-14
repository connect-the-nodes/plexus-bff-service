package com.zynchub.digital.hubservice.app.feature.flags.mapper;

import com.zynchub.digital.hubservice.app.feature.flags.dto.FeatureDto;
import com.zynchub.digital.hubservice.app.feature.flags.model.FeatureFlag;
import java.util.List;
import org.springframework.stereotype.Component;

@Component
public class FeatureMapper {

    public List<FeatureFlag> mapFromDto(List<FeatureDto> featureDtos) {
        return featureDtos.stream().map(this::mapFeature).toList();
    }

    private FeatureFlag mapFeature(FeatureDto featureDto) {
        FeatureFlag feature = new FeatureFlag();
        feature.setName(featureDto.name());
        feature.setParent(featureDto.parent() == null || featureDto.parent().isBlank()
                ? featureDto.name()
                : featureDto.parent());
        feature.setEnabled(featureDto.enabled());
        feature.setDescription(featureDto.description());
        return feature;
    }
}
