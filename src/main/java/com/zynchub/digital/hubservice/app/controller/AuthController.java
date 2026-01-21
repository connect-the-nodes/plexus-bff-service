package com.zynchub.digital.hubservice.app.controller;

import com.zynchub.digital.hubservice.app.config.AuthCognitoProperties;
import jakarta.servlet.http.HttpServletRequest;
import java.net.URI;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.UUID;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/auth")
public class AuthController {

  private final AuthCognitoProperties.Properties properties;

  public AuthController(AuthCognitoProperties.Properties properties) {
    this.properties = properties;
  }

  @GetMapping("/login")
  public ResponseEntity<Void> login() {
    if (!properties.enabled()) {
      return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    String state = UUID.randomUUID().toString();
    String scope = String.join(" ", defaultScopes(properties.scopes()));

    String authorizeUrl =
        String.format(
            "%s/oauth2/authorize?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s",
            properties.domain(),
            urlEncode(properties.clientId()),
            urlEncode(scope),
            urlEncode(properties.redirectUri()),
            urlEncode(state));

    HttpHeaders headers = new HttpHeaders();
    headers.setLocation(URI.create(authorizeUrl));
    return new ResponseEntity<>(headers, HttpStatus.FOUND);
  }

  @GetMapping("/callback")
  public ResponseEntity<Void> callback(
      @RequestParam(name = "code", required = false) String code,
      @RequestParam(name = "state", required = false) String state,
      HttpServletRequest request) {
    if (!properties.enabled()) {
      return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }

    if (!StringUtils.hasText(code)) {
      return ResponseEntity.status(HttpStatus.BAD_REQUEST).build();
    }

    String redirect = properties.postLoginRedirectUri();
    String target =
        String.format(
            "%s?code=%s%s",
            redirect,
            urlEncode(code),
            StringUtils.hasText(state) ? "&state=" + urlEncode(state) : "");

    HttpHeaders headers = new HttpHeaders();
    headers.setLocation(URI.create(target));
    return new ResponseEntity<>(headers, HttpStatus.FOUND);
  }

  private static List<String> defaultScopes(List<String> scopes) {
    return (scopes == null || scopes.isEmpty()) ? List.of("openid", "email", "profile") : scopes;
  }

  private static String urlEncode(String value) {
    return URLEncoder.encode(value, StandardCharsets.UTF_8);
  }
}
