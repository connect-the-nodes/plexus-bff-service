package com.zynchub.digital.hubservice.app.controller;

import jakarta.servlet.http.HttpSession;
import org.springframework.context.annotation.Profile;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Profile({"dev", "local-redis"})
public class TestSessionController {

  @PostMapping("/_test/session")
  public ResponseEntity<Void> createSession(HttpSession session) {
    session.setAttribute("test-key", "test-value");
    return ResponseEntity.noContent().build();
  }
}
