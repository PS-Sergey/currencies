package com.polyakov.currencies.feign;

import com.polyakov.currencies.vo.OpenExchangeRatesVo;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.Map;

@FeignClient(value = "openExchangeRatesFeign", url = "${openexchangerates.url}")
public interface OpenExchangeRatesFeign {

    @GetMapping("/currencies.json")
    Map<String, String> getCurrencies();

    @GetMapping("/historical/{date}.json")
    OpenExchangeRatesVo geHistorical(
        @PathVariable("date") String date,
        @RequestParam("app_id") String appId,
        @RequestParam("base") String base,
        @RequestParam("symbols") String symbols);
}
