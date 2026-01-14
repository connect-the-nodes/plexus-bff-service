package com.zynchub.digital.hubservice.app.repository.impl;


import com.zynchub.digital.hubservice.app.domain.StatusInfo;
import com.zynchub.digital.hubservice.app.repository.StatusRepository;
import org.springframework.stereotype.Repository;

@Repository
public class InMemoryStatusRepository implements StatusRepository {

    @Override
    public StatusInfo fetchStatus() {
        return new StatusInfo("OK", "Service is running");
    }
}
