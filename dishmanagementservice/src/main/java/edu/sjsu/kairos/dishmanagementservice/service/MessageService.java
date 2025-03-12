package edu.sjsu.kairos.dishmanagementservice.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.MessageSource;
import org.springframework.context.i18n.LocaleContextHolder;
import org.springframework.stereotype.Service;

import java.util.Locale;

@Service
public class MessageService {

	@Autowired
	private MessageSource messageSource;
	
	public String getMessage(String messageKey, Object... params) {
		Locale locale = LocaleContextHolder.getLocale();
		return messageSource.getMessage(messageKey, params, locale);
	}
	
}
