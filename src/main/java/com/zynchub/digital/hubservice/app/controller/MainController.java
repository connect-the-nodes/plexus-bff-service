package com.zynchub.digital.hubservice.app.controller;

import com.zynchub.digital.hubservice.app.dto.StatusResponseDto;
import com.zynchub.digital.hubservice.app.mapper.StatusMapper;
import com.zynchub.digital.hubservice.app.service.StatusService;
import org.springframework.web.bind.annotation.GetMapping;

public class MainController {

    private final StatusService statusService;

    public MainController(StatusService statusService) {
        this.statusService = statusService;
    }

    @GetMapping("/status")
    public StatusResponseDto status() {
        return StatusMapper.toDto(statusService.getStatus());
    }

}
