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
    public int getKey(String currency) {
        OpenExchangeRatesModel response = openExchangeRatesFeign.getLatest(appId, currency);
        Double rate = response.getRates().get(currency);

        LocalDate yesterday = LocalDate.now().minusDays(1);
        String date = DateTimeFormatter.ofPattern("yyyy-MM-dd").format(yesterday);
        OpenExchangeRatesModel yesterdayResponse = openExchangeRatesFeign.geHistorical(date, appId, currency);
        Double yesterdayRate = yesterdayResponse.getRates().get(currency);
        Double result = 1 / rate - 1 / yesterdayRate;
        return Double.compare(1 / rate, 1 / yesterdayRate);
    }
}
