package com.zynchub.digital.hubservice;

import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import java.io.IOException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;
import redis.embedded.RedisServer;

@Component
@Profile("local-redis")
public class EmbeddedRedisServer {

  private static final Logger log = LoggerFactory.getLogger(EmbeddedRedisServer.class);

  private final RedisServer redisServer;
  private boolean started;

  public EmbeddedRedisServer(@Value("${spring.data.redis.port}") int port) throws IOException {
    redisServer = new RedisServer(port);
  }

  @PostConstruct
  public void postConstruct() throws IOException {
    try {
      redisServer.start();
      started = true;
    } catch (RuntimeException ex) {
      log.warn("Embedded Redis already running or failed to start, continuing", ex);
    }
  }

  @PreDestroy
  public void preDestroy() throws IOException {
    if (started) {
      redisServer.stop();
    }
  }
}
