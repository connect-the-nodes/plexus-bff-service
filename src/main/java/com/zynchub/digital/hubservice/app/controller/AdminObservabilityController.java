package com.zynchub.digital.hubservice.app.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.dto.AdminObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.ConnectorHealthListResponseDto;
import com.zynchub.digital.hubservice.app.dto.TenantHealthListResponseDto;
import com.zynchub.digital.hubservice.app.service.ObservabilityService;
import java.time.Instant;
import java.util.Set;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping(ApiPaths.API_V3 + "/admin/observability")
public class AdminObservabilityController {

    private static final Set<String> ADMIN_AUTHORITIES = Set.of(
            "ROLE_ADMIN",
            "ROLE_AXIS_ADMIN",
            "ROLE_SUPER_ADMIN",
            "SCOPE_admin",
            "SCOPE_axis.admin");

    private final ObservabilityService observabilityService;

    public AdminObservabilityController(ObservabilityService observabilityService) {
        this.observabilityService = observabilityService;
    }

    @GetMapping("/overview")
    public ApiResponse<AdminObservabilityOverviewResponseDto> getOverview(
            @RequestParam(value = "from", required = false) Instant from,
            @RequestParam(value = "to", required = false) Instant to,
            @RequestParam(value = "periodSeconds", defaultValue = "300") int periodSeconds,
            @RequestParam(value = "tenantId", required = false) String tenantId,
            @RequestParam(value = "connectorId", required = false) String connectorId,
            Authentication authentication) {
        ensureAdmin(authentication);
        return ApiResponse.ok(
                observabilityService.getAdminOverview(from, to, periodSeconds, tenantId, connectorId));
    }

    @GetMapping("/tenants")
    public ApiResponse<TenantHealthListResponseDto> getTenantHealth(
            @RequestParam(value = "from", required = false) Instant from,
            @RequestParam(value = "to", required = false) Instant to,
            @RequestParam(value = "status", required = false) String status,
            @RequestParam(value = "page", defaultValue = "1") int page,
            @RequestParam(value = "size", defaultValue = "50") int size,
            Authentication authentication) {
        ensureAdmin(authentication);
        return ApiResponse.ok(observabilityService.getAdminTenantHealth(from, to, status, page, size));
    }

    @GetMapping("/connectors")
    public ApiResponse<ConnectorHealthListResponseDto> getConnectorHealth(
            @RequestParam(value = "from", required = false) Instant from,
            @RequestParam(value = "to", required = false) Instant to,
            @RequestParam(value = "status", required = false) String status,
            @RequestParam(value = "page", defaultValue = "1") int page,
            @RequestParam(value = "size", defaultValue = "50") int size,
            Authentication authentication) {
        ensureAdmin(authentication);
        return ApiResponse.ok(observabilityService.getAdminConnectorHealth(from, to, status, page, size));
    }

    private void ensureAdmin(Authentication authentication) {
        if (authentication == null) {
            return;
        }
        boolean authorized = authentication.getAuthorities().stream()
                .map(grantedAuthority -> grantedAuthority == null ? null : grantedAuthority.getAuthority())
                .anyMatch(ADMIN_AUTHORITIES::contains);
        if (!authorized) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Admin role required");
        }
    }
}
