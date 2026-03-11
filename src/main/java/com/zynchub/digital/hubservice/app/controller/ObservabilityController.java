package com.zynchub.digital.hubservice.app.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.dto.ObservabilityInventoryItemResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.service.ObservabilityService;
import java.time.Instant;
import java.util.List;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping(ApiPaths.API_V3 + "/observability")
public class ObservabilityController {

    private final ObservabilityService observabilityService;

    public ObservabilityController(ObservabilityService observabilityService) {
        this.observabilityService = observabilityService;
    }

    @GetMapping("/overview")
    public ApiResponse<ObservabilityOverviewResponseDto> getOverview(
            @RequestParam(value = "from", required = false) Instant from,
            @RequestParam(value = "to", required = false) Instant to,
            @RequestParam(value = "periodSeconds", defaultValue = "300") int periodSeconds,
            @RequestParam(value = "tenantId", required = false) String tenantId,
            Authentication authentication) {
        String resolvedTenantId = resolveTenantId(tenantId, authentication);
        return ApiResponse.ok(observabilityService.getOverview(from, to, periodSeconds, resolvedTenantId));
    }

    @GetMapping("/integrations/inventory")
    public ApiResponse<List<ObservabilityInventoryItemResponseDto>> getInventory(
            @RequestParam(value = "from", required = false) Instant from,
            @RequestParam(value = "to", required = false) Instant to,
            @RequestParam(value = "tenantId", required = false) String tenantId,
            Authentication authentication) {
        String resolvedTenantId = resolveTenantId(tenantId, authentication);
        return ApiResponse.ok(observabilityService.getInventory(from, to, resolvedTenantId));
    }

    private String resolveTenantId(String tenantIdParam, Authentication authentication) {
        if (authentication instanceof JwtAuthenticationToken jwtAuthenticationToken) {
            String jwtTenantId = getTenantFromJwt(jwtAuthenticationToken);
            if (isBlank(jwtTenantId)) {
                throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Tenant claim missing in JWT");
            }
            if (!isBlank(tenantIdParam) && !jwtTenantId.equals(tenantIdParam)) {
                throw new ResponseStatusException(HttpStatus.FORBIDDEN, "Tenant mismatch");
            }
            return jwtTenantId;
        }

        if (!isBlank(tenantIdParam)) {
            return tenantIdParam;
        }

        return "default";
    }

    private String getTenantFromJwt(JwtAuthenticationToken token) {
        Object tenantId = token.getToken().getClaims().get("tenantId");
        if (tenantId instanceof String tenantIdValue && !tenantIdValue.isBlank()) {
            return tenantIdValue;
        }

        Object customTenantId = token.getToken().getClaims().get("custom:tenantId");
        if (customTenantId instanceof String customTenantIdValue && !customTenantIdValue.isBlank()) {
            return customTenantIdValue;
        }

        Object altTenantId = token.getToken().getClaims().get("tenant_id");
        if (altTenantId instanceof String altTenantIdValue && !altTenantIdValue.isBlank()) {
            return altTenantIdValue;
        }

        return null;
    }

    private boolean isBlank(String value) {
        return value == null || value.isBlank();
    }
}
