package com.zynchub.digital.hubservice.app.dto;

import java.util.ArrayList;
import java.util.List;

public class ConnectorHealthListResponseDto {
    private long total;
    private List<ConnectorHealthSummaryResponseDto> items = new ArrayList<>();

    public ConnectorHealthListResponseDto() {
    }

    public long getTotal() {
        return total;
    }

    public void setTotal(long total) {
        this.total = total;
    }

    public List<ConnectorHealthSummaryResponseDto> getItems() {
        return items;
    }

    public void setItems(List<ConnectorHealthSummaryResponseDto> items) {
        this.items = items;
    }
}
