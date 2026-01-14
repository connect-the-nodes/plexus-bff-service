package com.zynchub.digital.hubservice.app.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI openAPI() {
        return new OpenAPI()
                .info(new Info()
                        .title("ZyncHub Digital Hub Service API")
                        .version("v3")
                        .description("REST endpoints for Hub Service (Connectors, Services, Transactions, etc.)"));
    }
}

