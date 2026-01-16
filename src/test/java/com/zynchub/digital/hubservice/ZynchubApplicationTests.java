package com.zynchub.digital.hubservice;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;

@SpringBootTest
@Import(com.zynchub.digital.hubservice.app.security.TestJwtDecoderConfig.class)
class ZynchubApplicationTests {

	@Test
	void contextLoads() {
	}

}
