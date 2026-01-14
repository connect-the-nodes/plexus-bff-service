package com.zynchub.digital.hubservice.app.repository;


import com.zynchub.digital.hubservice.app.domain.StatusInfo;

public interface StatusRepository {
    StatusInfo fetchStatus();
}
