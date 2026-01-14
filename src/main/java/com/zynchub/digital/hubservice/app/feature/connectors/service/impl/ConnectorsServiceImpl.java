package com.zynchub.digital.hubservice.app.feature.connectors.service.impl;

import com.zynchub.digital.hubservice.app.domain.StatusInfo;
import com.zynchub.digital.hubservice.app.feature.connectors.domain.ConnectorsInfo;
import com.zynchub.digital.hubservice.app.feature.connectors.dto.ConnectorsResponseDto;
import com.zynchub.digital.hubservice.app.feature.connectors.mapper.ConnectorsMapper;
import com.zynchub.digital.hubservice.app.feature.connectors.repository.ConnectorsRepository;
import com.zynchub.digital.hubservice.app.feature.connectors.service.ConnectorsService;
import com.zynchub.digital.hubservice.app.feature.status.dto.FeatureResponseDto;
import com.zynchub.digital.hubservice.app.repository.StatusRepository;
import org.springframework.stereotype.Service;

@Service
public class ConnectorsServiceImpl implements ConnectorsService {

    private final ConnectorsRepository connectorsRepository;
    private final ConnectorsMapper connectorsMapper;

    public ConnectorsServiceImpl(ConnectorsRepository connectorsRepository) {
        this.connectorsRepository = connectorsRepository;
        this.connectorsMapper=new ConnectorsMapper();
    }

    /**
     * @return
     */
    @Override
    public ConnectorsResponseDto getConnectors() {
        ConnectorsInfo connectorsInfo =  connectorsRepository.fetchConnectors();
        return connectorsMapper.mapConnectorsData(connectorsInfo);
    }



}

