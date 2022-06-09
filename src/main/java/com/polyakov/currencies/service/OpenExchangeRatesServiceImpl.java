package com.polyakov.currencies.service;

import com.polyakov.currencies.feign.OpenExchangeRatesFeign;
import com.polyakov.currencies.model.OpenExchangeRatesModel;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Service
public class OpenExchangeRatesServiceImpl implements OpenExchangeRatesService{

    private OpenExchangeRatesFeign openExchangeRatesFeign;

    @Autowired
    public OpenExchangeRatesServiceImpl(OpenExchangeRatesFeign openExchangeRatesFeign) {
        this.openExchangeRatesFeign = openExchangeRatesFeign;
    }

    @Override
    public List<String> getCurrencies(String appId) {
        OpenExchangeRatesModel exchangeModel = openExchangeRatesFeign.getLatest(appId);
        return exchangeModel
                .getRates()
                .entrySet()
                .stream()
                .map(Map.Entry::getKey)
                .collect(Collectors.toList());
    }
}
