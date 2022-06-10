package com.polyakov.currencies.controller;

import com.polyakov.currencies.service.OpenExchangeRatesService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api")
public class CurrencyController {

    private final OpenExchangeRatesService openExchangeRatesService;

    @Autowired
    public CurrencyController(OpenExchangeRatesService openExchangeRatesService) {
        this.openExchangeRatesService = openExchangeRatesService;
    }

    @GetMapping("/currencies")
    public Map<String, String> getCurrencies() {
        return openExchangeRatesService.getCurrencies();
    }

    @GetMapping("/rate/{currency}")
    public String getRate(@PathVariable("currency") String currency) {
        int key = openExchangeRatesService.getKey(currency);
        return Integer.toString(key);
    }
}
