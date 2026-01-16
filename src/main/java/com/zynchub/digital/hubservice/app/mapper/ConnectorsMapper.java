package com.zynchub.digital.hubservice.app.mapper;

import com.zynchub.digital.hubservice.app.model.ConnectorsInfo;
import com.zynchub.digital.hubservice.app.dto.ConnectorsResponseDto;

public class ConnectorsMapper {

    public ConnectorsResponseDto mapConnectorsData(ConnectorsInfo connectorsInfo){
        return new ConnectorsResponseDto(connectorsInfo.getStatus(), connectorsInfo.getMessage());
    }

}
