package com.zynchub.digital.hubservice.app.security;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.test.util.ReflectionTestUtils;

class JwtAuthConverterTest {

  @Test
  void convert_adds_group_and_scope_authorities() {
    JwtAuthConverter converter = new JwtAuthConverter();
    ReflectionTestUtils.setField(converter, "authoritiesClaim", "cognito:groups");
    ReflectionTestUtils.setField(converter, "rolePrefix", "ROLE_");

    Jwt jwt =
        Jwt.withTokenValue("token")
            .header("alg", "none")
            .claim("cognito:groups", List.of("ADMIN", "USER"))
            .claim("scope", List.of("read", "write"))
            .build();

    List<String> authorities =
        converter.convert(jwt).getAuthorities().stream()
            .map(GrantedAuthority::getAuthority)
            .sorted()
            .toList();

    assertThat(authorities)
        .containsExactly("ROLE_ADMIN", "ROLE_USER", "SCOPE_read", "SCOPE_write");
  }

  @Test
  void convert_handles_null_groups_and_scope_list() {
    JwtAuthConverter converter = new JwtAuthConverter();
    ReflectionTestUtils.setField(converter, "authoritiesClaim", "cognito:groups");
    ReflectionTestUtils.setField(converter, "rolePrefix", "ROLE_");

    Jwt jwt =
        Jwt.withTokenValue("token")
            .header("alg", "none")
            .claim("scope", List.of("profile"))
            .build();

    List<String> authorities =
        converter.convert(jwt).getAuthorities().stream()
            .map(GrantedAuthority::getAuthority)
            .sorted()
            .toList();

    assertThat(authorities).containsExactly("SCOPE_profile");
  }
}
