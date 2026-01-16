package com.zynchub.digital.hubservice.app.security;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.web.servlet.MockMvc;

@SpringBootTest(properties = "security.enabled=false")
@AutoConfigureMockMvc
class SecurityDisabledIT {

  @Autowired private MockMvc mockMvc;

  @Test
  void status_endpoint_is_open_when_security_disabled() throws Exception {
    mockMvc.perform(get(ApiPaths.API_V2 + "/status")).andExpect(status().isOk());
  }
}
