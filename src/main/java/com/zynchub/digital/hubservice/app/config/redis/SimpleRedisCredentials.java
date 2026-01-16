package com.zynchub.digital.hubservice.app.config.redis;

import io.lettuce.core.RedisCredentials;

class SimpleRedisCredentials implements RedisCredentials {

  private final String username;
  private final char[] password;

  SimpleRedisCredentials(String username, char[] password) {
    this.username = username;
    this.password = password;
  }

  @Override
  public String getUsername() {
    return username;
  }

  @Override
  public char[] getPassword() {
    return password;
  }

  @Override
  public boolean hasUsername() {
    return username != null && !username.isBlank();
  }

  @Override
  public boolean hasPassword() {
    return password != null && password.length > 0;
  }
}
