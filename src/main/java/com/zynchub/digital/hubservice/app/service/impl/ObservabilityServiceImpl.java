package com.zynchub.digital.hubservice.app.service.impl;

import com.zynchub.digital.hubservice.app.config.ObservabilityClientProperties;
import com.zynchub.digital.hubservice.app.dto.AdminObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.ConnectorHealthListResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityInventoryItemResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.TenantHealthListResponseDto;
import com.zynchub.digital.hubservice.app.service.ObservabilityService;
import java.time.Instant;
import java.util.List;
import java.util.Optional;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

@Service
public class ObservabilityServiceImpl implements ObservabilityService {

    private final RestClient restClient;

    public ObservabilityServiceImpl(
            RestClient.Builder restClientBuilder,
            ObservabilityClientProperties.Properties properties) {
        this.restClient = restClientBuilder.baseUrl(properties.baseUrl()).build();
    }

    @Override
    public ObservabilityOverviewResponseDto getOverview(
            Instant from, Instant to, int periodSeconds, String tenantId) {
        return restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/v1/observability/overview")
                        .queryParamIfPresent("from", Optional.ofNullable(from))
                        .queryParamIfPresent("to", Optional.ofNullable(to))
                        .queryParamIfPresent("tenantId", Optional.ofNullable(tenantId))
                        .queryParam("periodSeconds", periodSeconds)
                        .build())
                .retrieve()
                .body(ObservabilityOverviewResponseDto.class);
    }

    @Override
    public List<ObservabilityInventoryItemResponseDto> getInventory(
            Instant from, Instant to, String tenantId) {
        List<ObservabilityInventoryItemResponseDto> inventory = restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/v1/observability/integrations/inventory")
                        .queryParamIfPresent("from", Optional.ofNullable(from))
                        .queryParamIfPresent("to", Optional.ofNullable(to))
                        .queryParamIfPresent("tenantId", Optional.ofNullable(tenantId))
                        .build())
                .retrieve()
                .body(new ParameterizedTypeReference<List<ObservabilityInventoryItemResponseDto>>() {
                });

        return inventory == null ? List.of() : inventory;
    }

    @Override
    public AdminObservabilityOverviewResponseDto getAdminOverview(
            Instant from,
            Instant to,
            int periodSeconds,
            String tenantId,
            String connectorId) {
        return restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/v1/admin/observability/overview")
                        .queryParamIfPresent("from", Optional.ofNullable(from))
                        .queryParamIfPresent("to", Optional.ofNullable(to))
                        .queryParam("periodSeconds", periodSeconds)
                        .queryParamIfPresent("tenantId", Optional.ofNullable(tenantId))
                        .queryParamIfPresent("connectorId", Optional.ofNullable(connectorId))
                        .build())
                .retrieve()
                .body(AdminObservabilityOverviewResponseDto.class);
    }

    @Override
    public TenantHealthListResponseDto getAdminTenantHealth(
            Instant from, Instant to, String status, int page, int size) {
        TenantHealthListResponseDto response = restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/v1/admin/observability/tenants")
                        .queryParamIfPresent("from", Optional.ofNullable(from))
                        .queryParamIfPresent("to", Optional.ofNullable(to))
                        .queryParamIfPresent("status", Optional.ofNullable(status))
                        .queryParam("page", page)
                        .queryParam("size", size)
                        .build())
                .retrieve()
                .body(TenantHealthListResponseDto.class);
        return response == null ? new TenantHealthListResponseDto() : response;
    }

    @Override
    public ConnectorHealthListResponseDto getAdminConnectorHealth(
            Instant from, Instant to, String status, int page, int size) {
        ConnectorHealthListResponseDto response = restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/v1/admin/observability/connectors")
                        .queryParamIfPresent("from", Optional.ofNullable(from))
                        .queryParamIfPresent("to", Optional.ofNullable(to))
                        .queryParamIfPresent("status", Optional.ofNullable(status))
                        .queryParam("page", page)
                        .queryParam("size", size)
                        .build())
                .retrieve()
                .body(ConnectorHealthListResponseDto.class);
        return response == null ? new ConnectorHealthListResponseDto() : response;
    }
}
