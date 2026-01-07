package com.zynchub.digital.hubservice.controller;


import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1")
@RequiredArgsConstructor
public class SessionController {


    // New simple test endpoint
    @GetMapping("/hello")
    public ResponseEntity<String> helloZynchHub() {
        return ResponseEntity.ok("Hello ZynchHub");
    }
}
