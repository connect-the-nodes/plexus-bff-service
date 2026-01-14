package com.zynchub.digital.hubservice.app.controller;


import com.zynchub.digital.hubservice.app.dto.StatusResponseDto;
import com.zynchub.digital.hubservice.app.mapper.StatusMapper;
import com.zynchub.digital.hubservice.app.service.StatusService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1")
@Tag(name = "Status", description = "Service health and status endpoints")
public class StatusController {

    private final StatusService statusService;

    public StatusController(StatusService statusService) {
        this.statusService = statusService;
    }

    @GetMapping("/status")
    @Operation(summary = "Get service status", description = "Returns current health/status of the application")
    public StatusResponseDto status() {
        return StatusMapper.toDto(statusService.getStatus());
    }
}
