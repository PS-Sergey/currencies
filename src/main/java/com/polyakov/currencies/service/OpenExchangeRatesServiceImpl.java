package com.polyakov.currencies.service;

import com.polyakov.currencies.exception.CurrencyRateNotFoundException;
import com.polyakov.currencies.feign.OpenExchangeRatesFeign;
import com.polyakov.currencies.vo.OpenExchangeRatesVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.time.LocalDate;
import java.util.Map;
import java.util.Objects;

@Service
public class OpenExchangeRatesServiceImpl implements OpenExchangeRatesService{

    @Value("${openexchangerates.app.id}")
    private String appId;

    @Value("${openexchangerates.base}")
    private String base;

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
    public int compareRate(String code) {
        LocalDate today = LocalDate.now();
        Double todayRate = getRateByDate(today.toString(), code);
        Double yesterdayRate = getRateByDate(today.minusDays(1).toString(), code);
        return Double.compare(todayRate, yesterdayRate);
    }

    private Double getRateByDate(String date, String code) {
        OpenExchangeRatesVo response = openExchangeRatesFeign.geHistorical(date, appId, base, code);
        Double rate = response.getRates().get(code);
        if (Objects.isNull(rate)) {
            throw new CurrencyRateNotFoundException(String.format("Currency rate for %s not found", code));
        }
        return 1 / rate;
    }
}
