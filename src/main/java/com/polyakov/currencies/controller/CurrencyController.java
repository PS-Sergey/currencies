package com.polyakov.currencies.controller;

import com.polyakov.currencies.service.OpenExchangeRatesService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api")
public class CurrencyController {

    @Value("${openexchangerates.app.id}")
    private String appId;

    private final OpenExchangeRatesService openExchangeRatesService;

    @Autowired
    public CurrencyController(OpenExchangeRatesService openExchangeRatesService) {
        this.openExchangeRatesService = openExchangeRatesService;
    }

    @GetMapping("/currencies")
    public List<String> getCurrencies() {
        return openExchangeRatesService.getCurrencies(appId);
    }
}
