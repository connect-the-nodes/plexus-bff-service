package com.zynchub.digital.hubservice.app.mapper;


import com.zynchub.digital.hubservice.app.domain.StatusInfo;
import com.zynchub.digital.hubservice.app.dto.StatusResponseDto;
import com.zynchub.digital.hubservice.app.util.TimeUtil;

public final class StatusMapper {
    private StatusMapper() {}

    public static StatusResponseDto toDto(StatusInfo info) {
        return new StatusResponseDto(
                info.getStatus(),
                info.getMessage(),
                TimeUtil.nowIso()
        );
    }
}

