package com.zynchub.digital.hubservice.app.feature.connectors.repository;

import com.zynchub.digital.hubservice.app.feature.connectors.domain.ConnectorsInfo;

public interface ConnectorsRepository {

    ConnectorsInfo fetchConnectors();

}
