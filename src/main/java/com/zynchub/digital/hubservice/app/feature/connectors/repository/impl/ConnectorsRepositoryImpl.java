package com.zynchub.digital.hubservice.app.feature.connectors.repository.impl;

import com.fasterxml.jackson.databind.util.JSONPObject;
import com.nimbusds.jose.shaded.gson.JsonObject;
import com.zynchub.digital.hubservice.app.feature.connectors.domain.ConnectorsInfo;
import com.zynchub.digital.hubservice.app.feature.connectors.repository.ConnectorsRepository;
import org.springframework.stereotype.Repository;

@Repository
public class ConnectorsRepositoryImpl implements ConnectorsRepository {

    /**
     * @implNote Retrieves the connectors details from the database
     * and sends it back to the product for displaying it the list of
     * available Connectors dashboard
     * @return {@link ConnectorsInfo
     * }
     */
    @Override
    public ConnectorsInfo fetchConnectors() {


        return new ConnectorsInfo("OK", "Connectors Service is running");
    }
}
