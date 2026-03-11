package com.zynchub.digital.hubservice.app.service;

import com.zynchub.digital.hubservice.app.dto.AdminObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.ConnectorHealthListResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityInventoryItemResponseDto;
import com.zynchub.digital.hubservice.app.dto.ObservabilityOverviewResponseDto;
import com.zynchub.digital.hubservice.app.dto.TenantHealthListResponseDto;
import java.time.Instant;
import java.util.List;

public interface ObservabilityService {

    ObservabilityOverviewResponseDto getOverview(Instant from, Instant to, int periodSeconds, String tenantId);

    List<ObservabilityInventoryItemResponseDto> getInventory(Instant from, Instant to, String tenantId);

    AdminObservabilityOverviewResponseDto getAdminOverview(
            Instant from,
            Instant to,
            int periodSeconds,
            String tenantId,
            String connectorId);

    TenantHealthListResponseDto getAdminTenantHealth(
            Instant from,
            Instant to,
            String status,
            int page,
            int size);

    ConnectorHealthListResponseDto getAdminConnectorHealth(
            Instant from,
            Instant to,
            String status,
            int page,
            int size);
}
