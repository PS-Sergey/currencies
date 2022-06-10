package com.polyakov.currencies.service;

import com.polyakov.currencies.feign.OpenExchangeRatesFeign;
import com.polyakov.currencies.model.OpenExchangeRatesModel;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.format.datetime.DateFormatter;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Service
public class OpenExchangeRatesServiceImpl implements OpenExchangeRatesService{

    @Value("${openexchangerates.app.id}")
    private String appId;
    @Value("${openexchangerates.base}")
    private String currency;

    private final OpenExchangeRatesFeign openExchangeRatesFeign;

    @Autowired
    public OpenExchangeRatesServiceImpl(OpenExchangeRatesFeign openExchangeRatesFeign) {
        this.openExchangeRatesFeign = openExchangeRatesFeign;
    }

    @Override
    public Map<String, String> getCurrencies() {
        return openExchangeRatesFeign.getCurrencies();
    }

    @Override
    public int compareRate() {
        LocalDate today = LocalDate.now();
        Double todayRate = getRateByDate(DateTimeFormatter.ofPattern("yyyy-MM-dd").format(today));
        Double yesterdayRate = getRateByDate(DateTimeFormatter.ofPattern("yyyy-MM-dd").format(today.minusDays(1)));
        return Double.compare(todayRate, yesterdayRate);
    }

    private Double getRateByDate(String date) {
        OpenExchangeRatesModel response = openExchangeRatesFeign.geHistorical(date, appId, currency);
        return 1 / response.getRates().get(currency);
    }
}
