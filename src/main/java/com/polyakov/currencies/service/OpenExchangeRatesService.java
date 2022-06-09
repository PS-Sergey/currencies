package com.polyakov.currencies.service;

import java.util.List;

public interface OpenExchangeRatesService {

    List<String> getCurrencies(String appId);
}
