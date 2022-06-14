package com.polyakov.currencies.exception;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.*;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ErrorMessage {
    String error;
    @JsonProperty("error_description")
    String errorDescription;
}
