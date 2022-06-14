package com.polyakov.currencies.service;

import java.util.Map;

public interface OpenExchangeRatesService {

    Map<String, String> getCurrencies();
    int compareRate();
}
