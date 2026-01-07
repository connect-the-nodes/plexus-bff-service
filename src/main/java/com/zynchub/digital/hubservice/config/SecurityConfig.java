package com.zynchub.digital.hubservice.config;


import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
public class SecurityConfig {

    @Bean
    SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        http
                .csrf(csrf -> csrf.disable())
                .authorizeHttpRequests(auth -> auth
                        .requestMatchers(
                                "/swagger-ui/**",
                                "/v3/api-docs/**",
                                "/api/v1/hello"   // <-- public
                        ).permitAll()
                        .anyRequest().authenticated()
                )
                .formLogin(form -> form
                        .permitAll()  // allows login page for everyone
                )
                .logout(logout -> logout
                        .permitAll()
                );

        return http.build();
    }
}
