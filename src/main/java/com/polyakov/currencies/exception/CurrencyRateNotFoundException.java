package com.polyakov.currencies.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.BAD_GATEWAY)
public class CurrencyRateNotFoundException extends RuntimeException {

    public CurrencyRateNotFoundException(String message) {
        super(message);
    }
}
