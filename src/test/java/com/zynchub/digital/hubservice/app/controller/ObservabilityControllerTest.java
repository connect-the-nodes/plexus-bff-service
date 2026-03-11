package com.zynchub.digital.hubservice.app.controller;

import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.zynchub.digital.hubservice.app.dto.ObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.service.ObservabilityService;
import java.time.Instant;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.web.server.ResponseStatusException;

class ObservabilityControllerTest {

  @Test
  void usesTenantFromJwtClaimForOverview() {
    ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
    ObservabilityController controller = new ObservabilityController(observabilityService);
    ObservabilityOverviewResponseDto response = new ObservabilityOverviewResponseDto();

    when(observabilityService.getOverview(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq(300), eq("tenant-jwt")))
        .thenReturn(response);

    controller.getOverview(null, null, 300, null, jwtAuth("tenant-jwt"));

    verify(observabilityService)
        .getOverview(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq(300), eq("tenant-jwt"));
  }

  @Test
  void rejectsTenantMismatchBetweenQueryAndJwt() {
    ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
    ObservabilityController controller = new ObservabilityController(observabilityService);

    assertThatThrownBy(
            () ->
                controller.getOverview(
                    Instant.now().minusSeconds(60),
                    Instant.now(),
                    300,
                    "tenant-query",
                    jwtAuth("tenant-jwt")))
        .isInstanceOf(ResponseStatusException.class)
        .satisfies(
            ex ->
                org.assertj.core.api.Assertions.assertThat(((ResponseStatusException) ex).getStatusCode())
                    .isEqualTo(HttpStatus.FORBIDDEN));
  }

  @Test
  void fallsBackToDefaultTenantInLocalModeWhenJwtMissing() {
    ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
    ObservabilityController controller = new ObservabilityController(observabilityService);

    when(observabilityService.getInventory(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq("default")))
        .thenReturn(List.of());

    controller.getInventory(null, null, null, null);

    verify(observabilityService)
        .getInventory(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq("default"));
  }

  @Test
  void allowsTenantIdQueryInLocalModeWhenAuthMissing() {
    ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
    ObservabilityController controller = new ObservabilityController(observabilityService);

    when(observabilityService.getInventory(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq("tenant-local")))
        .thenReturn(List.of());

    controller.getInventory(null, null, "tenant-local", null);

    verify(observabilityService)
        .getInventory(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq("tenant-local"));
  }

  @Test
  void usesTenantFromCognitoGroupWhenTenantClaimMissing() {
    ObservabilityService observabilityService = org.mockito.Mockito.mock(ObservabilityService.class);
    ObservabilityController controller = new ObservabilityController(observabilityService);
    ObservabilityOverviewResponseDto response = new ObservabilityOverviewResponseDto();

    when(observabilityService.getOverview(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq(300), eq("demo-a")))
        .thenReturn(response);

    Jwt jwt =
        Jwt.withTokenValue("token")
            .header("alg", "none")
            .claim("cognito:groups", List.of("TENANT_demo-a"))
            .build();
    JwtAuthenticationToken auth = new JwtAuthenticationToken(jwt, List.of(new SimpleGrantedAuthority("ROLE_ADMIN")));

    controller.getOverview(null, null, 300, null, auth);

    verify(observabilityService)
        .getOverview(org.mockito.ArgumentMatchers.isNull(), org.mockito.ArgumentMatchers.isNull(), eq(300), eq("demo-a"));
  }

  private JwtAuthenticationToken jwtAuth(String tenantId) {
    Jwt jwt =
        Jwt.withTokenValue("token")
            .header("alg", "none")
            .claim("tenantId", tenantId)
            .claim("cognito:groups", List.of("ADMIN"))
            .build();
    return new JwtAuthenticationToken(jwt, List.of(new SimpleGrantedAuthority("ROLE_ADMIN")));
  }
}
