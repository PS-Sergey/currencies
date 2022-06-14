package com.polyakov.currencies.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class CurrencyRateNotFoundException extends RuntimeException {

    public CurrencyRateNotFoundException(String message) {
        super(message);
    }
}
