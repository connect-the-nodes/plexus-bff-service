package com.zynchub.digital.hubservice.app.config;

import java.util.Map;
import org.slf4j.MDC;
import org.springframework.core.task.TaskDecorator;

public class MdcTaskDecorator implements TaskDecorator {

  @Override
  public Runnable decorate(Runnable runnable) {
    Map<String, String> contextMap = MDC.getCopyOfContextMap();
    return () -> {
      Map<String, String> previous = MDC.getCopyOfContextMap();
      try {
        if (contextMap == null) {
          MDC.clear();
        } else {
          MDC.setContextMap(contextMap);
        }
        runnable.run();
      } finally {
        if (previous == null) {
          MDC.clear();
        } else {
          MDC.setContextMap(previous);
        }
      }
    };
  }
}
