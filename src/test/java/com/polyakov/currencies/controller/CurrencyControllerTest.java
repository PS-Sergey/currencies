package com.polyakov.currencies.controller;

import com.polyakov.currencies.service.GiphyService;
import com.polyakov.currencies.service.OpenExchangeRatesService;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import java.util.HashMap;
import java.util.Map;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.content;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(CurrencyController.class)
class CurrencyControllerTest {

    @Value("${giphy.rich}")
    private String rich;

    @MockBean
    private OpenExchangeRatesService openExchangeRatesService;

    @MockBean
    private GiphyService giphyService;

    @Autowired
    private MockMvc mockMvc;

    @Test
    void whenGetCurrenciesThenReturnCurrencies() throws Exception {
        Map<String, String> responseMap = new HashMap<>();
        responseMap.put("RUB", "Russian Ruble");
        String expectingResult = "{\"RUB\":\"Russian Ruble\"}";
        Mockito.when(openExchangeRatesService.getCurrencies())
            .thenReturn(responseMap);
        mockMvc.perform(get("/api/currencies"))
            .andDo(print())
            .andExpect(content().json(expectingResult));
    }

    @Test
    void whenGetRateThenReturnStatusOkAndContentType() throws Exception {
        Mockito.when(openExchangeRatesService.compareRate("RUB"))
            .thenReturn(1);
        Mockito.when(giphyService.getGif(rich))
            .thenReturn(new byte[1]);
        mockMvc.perform(get("/api/rate"))
            .andDo(print())
            .andExpect(content().contentType(MediaType.IMAGE_GIF))
            .andExpect(status().isOk());
    }
}