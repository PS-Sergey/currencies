package com.polyakov.currencies.feign;

import com.polyakov.currencies.model.OpenExchangeRatesModel;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.List;
import java.util.Map;

@FeignClient(value = "OpenExchangeRatesFeign", url = "${openexchangerates.url}")
public interface OpenExchangeRatesFeign {

    @GetMapping("/currencies.json")
    Map<String, String> getCurrencies();

    @GetMapping("/latest.json")
    OpenExchangeRatesModel getLatest(
        @RequestParam("app_id") String appId,
        @RequestParam("symbols") String symbols);

    @GetMapping("/historical/{date}.json")
    OpenExchangeRatesModel geHistorical(
        @PathVariable("date") String date,
        @RequestParam("app_id") String appId,
        @RequestParam("symbols") String symbols);
}
