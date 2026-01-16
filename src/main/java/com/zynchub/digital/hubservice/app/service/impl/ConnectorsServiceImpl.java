package com.zynchub.digital.hubservice.app.service.impl;

import com.zynchub.digital.hubservice.app.model.ConnectorsInfo;
import com.zynchub.digital.hubservice.app.dto.ConnectorsResponseDto;
import com.zynchub.digital.hubservice.app.mapper.ConnectorsMapper;
import com.zynchub.digital.hubservice.app.repository.ConnectorsRepository;
import com.zynchub.digital.hubservice.app.service.ConnectorsService;
import com.zynchub.digital.hubservice.app.feature.FeatureIds;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import org.springframework.stereotype.Service;
import org.springframework.security.access.AccessDeniedException;

@Service
public class ConnectorsServiceImpl implements ConnectorsService {

    private final ConnectorsRepository connectorsRepository;
    private final ConnectorsMapper connectorsMapper;
    private final FeaturesRetriever featuresRetriever;

    public ConnectorsServiceImpl(ConnectorsRepository connectorsRepository,
                                 FeaturesRetriever featuresRetriever) {
        this.connectorsRepository = connectorsRepository;
        this.connectorsMapper = new ConnectorsMapper();
        this.featuresRetriever = featuresRetriever;
    }

    /**
     * @return
     */
    @Override
    public ConnectorsResponseDto getConnectors() {
        if (!featuresRetriever.isActive(FeatureIds.CONNECTORS)) {
            throw new AccessDeniedException("Feature " + FeatureIds.CONNECTORS + " is not enabled!");
        }
        ConnectorsInfo connectorsInfo =  connectorsRepository.fetchConnectors();
        return connectorsMapper.mapConnectorsData(connectorsInfo);
    }



}

