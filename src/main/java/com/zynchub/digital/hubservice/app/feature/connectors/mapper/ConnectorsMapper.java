package com.zynchub.digital.hubservice.app.feature.connectors.mapper;

import com.zynchub.digital.hubservice.app.feature.connectors.domain.ConnectorsInfo;
import com.zynchub.digital.hubservice.app.feature.connectors.dto.ConnectorsResponseDto;

public class ConnectorsMapper {

    public ConnectorsResponseDto mapConnectorsData(ConnectorsInfo connectorsInfo){
        return new ConnectorsResponseDto(connectorsInfo.getStatus(), connectorsInfo.getMessage());
    }

}
