package com.polyakov.currencies.service;

import java.util.List;
import java.util.Map;

public interface OpenExchangeRatesService {

    Map<String, String> getCurrencies();
    int getKey(String currency);
}
