package com.zynchub.digital.hubservice.app.security;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Objects;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.convert.converter.Converter;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Component;

@Component
public class JwtAuthConverter implements Converter<Jwt, AbstractAuthenticationToken> {

  @Value("${security.jwt.authorities-claim:cognito:groups}")
  private String authoritiesClaim;

  @Value("${security.jwt.role-prefix:ROLE_}")
  private String rolePrefix;

  @Override
  public AbstractAuthenticationToken convert(Jwt jwt) {
    List<GrantedAuthority> authorities = new ArrayList<>();

    Object groups = jwt.getClaim(authoritiesClaim);
    if (groups instanceof Collection<?> groupCollection) {
      for (Object group : groupCollection) {
        if (group != null) {
          authorities.add(new SimpleGrantedAuthority(rolePrefix + group));
        }
      }
    }

    Object scopes = jwt.getClaim("scope");
    if (scopes instanceof Collection<?> scopeCollection) {
      for (Object scope : scopeCollection) {
        if (scope != null) {
          authorities.add(new SimpleGrantedAuthority("SCOPE_" + scope));
        }
      }
    } else if (scopes instanceof String scopeString) {
      for (String scope : scopeString.split(" ")) {
        if (!scope.isBlank()) {
          authorities.add(new SimpleGrantedAuthority("SCOPE_" + scope));
        }
      }
    }

    List<GrantedAuthority> filteredAuthorities =
        authorities.stream().filter(Objects::nonNull).toList();
    return new JwtAuthenticationToken(jwt, filteredAuthorities);
  }
}
