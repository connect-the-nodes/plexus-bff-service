package com.zynchub.digital.hubservice.app.controller;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.zynchub.digital.hubservice.app.dto.AdminObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.TenantHealthListResponseDto;
import com.zynchub.digital.hubservice.app.service.ObservabilityService;
import java.time.Instant;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.web.server.ResponseStatusException;

class AdminObservabilityControllerTest {

    @Test
    void allows_admin_role_for_overview_endpoint() {
        ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
        AdminObservabilityController controller = new AdminObservabilityController(observabilityService);
        AdminObservabilityOverviewResponseDto response = new AdminObservabilityOverviewResponseDto();

        when(observabilityService.getAdminOverview(
                org.mockito.ArgumentMatchers.isNull(),
                org.mockito.ArgumentMatchers.isNull(),
                eq(300),
                eq("tenant-a"),
                eq("dvla")))
            .thenReturn(response);

        controller.getOverview(null, null, 300, "tenant-a", "dvla", jwtAuth("ROLE_ADMIN"));

        verify(observabilityService).getAdminOverview(
                org.mockito.ArgumentMatchers.isNull(),
                org.mockito.ArgumentMatchers.isNull(),
                eq(300),
                eq("tenant-a"),
                eq("dvla"));
    }

    @Test
    void rejects_non_admin_access() {
        ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
        AdminObservabilityController controller = new AdminObservabilityController(observabilityService);

        assertThatThrownBy(() ->
                controller.getOverview(
                        Instant.now().minusSeconds(60),
                        Instant.now(),
                        300,
                        null,
                        null,
                        jwtAuth("ROLE_TENANT_VIEWER")))
            .isInstanceOf(ResponseStatusException.class)
            .satisfies(ex ->
                org.assertj.core.api.Assertions.assertThat(((ResponseStatusException) ex).getStatusCode())
                    .isEqualTo(HttpStatus.FORBIDDEN));
    }

    @Test
    void allows_local_no_auth_for_admin_tenant_health() {
        ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
        AdminObservabilityController controller = new AdminObservabilityController(observabilityService);
        TenantHealthListResponseDto response = new TenantHealthListResponseDto();

        when(observabilityService.getAdminTenantHealth(
                org.mockito.ArgumentMatchers.isNull(),
                org.mockito.ArgumentMatchers.isNull(),
                eq("critical"),
                eq(2),
                eq(25)))
            .thenReturn(response);

        controller.getTenantHealth(null, null, "critical", 2, 25, null);

        verify(observabilityService).getAdminTenantHealth(
                org.mockito.ArgumentMatchers.isNull(),
                org.mockito.ArgumentMatchers.isNull(),
                eq("critical"),
                eq(2),
                eq(25));
    }

    private JwtAuthenticationToken jwtAuth(String authority) {
        Jwt jwt =
            Jwt.withTokenValue("token")
                .header("alg", "none")
                .claim("tenantId", "tenant-jwt")
                .claim("cognito:groups", List.of("ADMIN"))
                .build();
        return new JwtAuthenticationToken(jwt, List.of(new SimpleGrantedAuthority(authority)));
    }
}
