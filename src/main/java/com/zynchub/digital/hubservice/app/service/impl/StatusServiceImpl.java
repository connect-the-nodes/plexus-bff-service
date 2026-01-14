package com.zynchub.digital.hubservice.app.service.impl;


import com.zynchub.digital.hubservice.app.domain.StatusInfo;
import com.zynchub.digital.hubservice.app.repository.StatusRepository;
import com.zynchub.digital.hubservice.app.service.StatusService;
import org.springframework.stereotype.Service;

@Service
public class StatusServiceImpl implements StatusService {

    private final StatusRepository statusRepository;

    public StatusServiceImpl(StatusRepository statusRepository) {
        this.statusRepository = statusRepository;
    }

    @Override
    public StatusInfo getStatus() {
        return statusRepository.fetchStatus();
    }
}
