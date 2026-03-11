package com.zynchub.digital.hubservice.app.dto;

import java.util.ArrayList;
import java.util.List;

public class TenantHealthListResponseDto {
    private long total;
    private List<TenantHealthSummaryResponseDto> items = new ArrayList<>();

    public TenantHealthListResponseDto() {
    }

    public long getTotal() {
        return total;
    }

    public void setTotal(long total) {
        this.total = total;
    }

    public List<TenantHealthSummaryResponseDto> getItems() {
        return items;
    }

    public void setItems(List<TenantHealthSummaryResponseDto> items) {
        this.items = items;
    }
}
