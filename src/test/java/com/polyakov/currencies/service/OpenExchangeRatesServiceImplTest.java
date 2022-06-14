package com.polyakov.currencies.service;

import com.polyakov.currencies.feign.OpenExchangeRatesFeign;
import com.polyakov.currencies.vo.OpenExchangeRatesVo;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;

import java.util.HashMap;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
class OpenExchangeRatesServiceImplTest {

    @Value("${openexchangerates.app.id}")
    private String appId;

    @Value("${openexchangerates.base}")
    private String base;

    private String code = "RUB";

    @MockBean
    private OpenExchangeRatesFeign openExchangeRatesFeign;

    @Autowired
    private OpenExchangeRatesServiceImpl openExchangeRatesService;

    @Test
    void whenGetCurrenciesThenReturnCurrencies() {
        Map<String, String> responseMap = new HashMap<>();
        responseMap.put("RUB", "Russian Ruble");
        Mockito.when(openExchangeRatesFeign.getCurrencies())
            .thenReturn(responseMap);

        Map<String, String> result = openExchangeRatesService.getCurrencies();
        assertEquals(responseMap, result);
    }

    @Test
    void whenCompareRateThenReturnNegative() {
        String currentDay = "2022-06-14";
        String oldDay = "2022-06-13";

        Mockito.when(openExchangeRatesFeign.geHistorical(currentDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(58.579, 1655218798));
        Mockito.when(openExchangeRatesFeign.geHistorical(oldDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(57.749, 1655164784));

        int result = openExchangeRatesService.compareRate(code);
        assertEquals(-1, result);
    }

    @Test
    void whenCompareRateThenReturnPositive() {
        String currentDay = "2022-06-14";
        String oldDay = "2022-06-13";

        Mockito.when(openExchangeRatesFeign.geHistorical(currentDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(57.749, 1655164784));
        Mockito.when(openExchangeRatesFeign.geHistorical(oldDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(58.579, 1655218798));

        int result = openExchangeRatesService.compareRate(code);
        assertEquals(1, result);
    }

    @Test
    void whenCompareRateThenReturnEquals() {
        String currentDay = "2022-06-14";
        String oldDay = "2022-06-13";

        Mockito.when(openExchangeRatesFeign.geHistorical(currentDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(57.749, 1655164784));
        Mockito.when(openExchangeRatesFeign.geHistorical(oldDay, appId, base, code))
            .thenReturn(createOpenExchangeRatesVo(57.749, 1655218798));

        int result = openExchangeRatesService.compareRate(code);
        assertEquals(0, result);
    }

    private OpenExchangeRatesVo createOpenExchangeRatesVo (Double rateValue, int imestamp) {
        Map<String, Double> rate = Map.of(code, rateValue);
        OpenExchangeRatesVo rateVo = new OpenExchangeRatesVo();
        rateVo.setDisclaimer("currentDayRateVo");
        rateVo.setLicense("https://openexchangerates.org/license");
        rateVo.setTimestamp(imestamp);
        rateVo.setBase("USD");
        rateVo.setRates(rate);
        return rateVo;
    }
}