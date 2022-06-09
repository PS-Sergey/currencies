package com.polyakov.currencies.feign;

import com.polyakov.currencies.model.OpenExchangeRatesModel;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(value = "OpenExchangeRatesFeign", url = "${openexchangerates.url}")
public interface OpenExchangeRatesFeign {

    @GetMapping("/latest.json")
    OpenExchangeRatesModel getLatest(@RequestParam("app_id") String appId);
}
