package com.zynchub.digital.hubservice.app.exception;

import com.zynchub.digital.hubservice.app.dto.ErrorResponseDto;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

public class ApiExceptionHandler {

    //@ExceptionHandler(NotFoundException.class)
    public ErrorResponseDto handleNotFound(NotFoundException ex) {
        return new ErrorResponseDto("NOT_FOUND", ex.getMessage());
    }

    //@ExceptionHandler(Exception.class)
    public ErrorResponseDto handleGeneric(Exception ex) {
        return new ErrorResponseDto("INTERNAL_ERROR", "Something went wrong");
    }
}
