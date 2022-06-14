package com.polyakov.currencies.feign;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

@FeignClient(name = "giphyFeign", url = "${giphy.url}")
public interface GiphyFeign {

    @GetMapping("/random")
    String getGif(
            @RequestParam("api_key") String apiKey,
            @RequestParam("tag") String tag
    );
}
