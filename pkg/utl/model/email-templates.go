package model

const (
	ActionTemplate = `<!-- Section-1 -->
	<table class="table_full editable-bg-color bg_color_e6e6e6 editable-bg-image" bgcolor="#e6e6e6" width="100%"
		align="center" mc:repeatable="castellab" mc:variant="Header" cellspacing="0" cellpadding="0" border="0">
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_303f9f" bgcolor="#303f9f" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="25"></td>
					</tr>
					<tr>
						<td>
							<!-- Inner container -->
							<table class="table1" width="520" align="center" border="0" cellspacing="0" cellpadding="0"
								style="margin: 0 auto;">
								<!-- horizontal gap -->
								<tr>
									<td height="40"></td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="20"></td>
								</tr>
	
								<tr>
									<td mc:edit="text001" align="center" class="text_color_ffffff"
										style="color: #ffffff; font-size: 30px; font-weight: 700; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text">
											<span class="text_container">
												<multiline>
													{{.Subject}}
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="30"></td>
								</tr>
							</table><!-- END inner container -->
						</td>
					</tr>
					<!-- padding-bottom -->
					<tr>
						<td height="60"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_ffffff" bgcolor="#ffffff" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="60"></td>
					</tr>
	
					<tr>
						<td>
							<!-- container_400 -->
							<table class="container_400" align="center" width="400" border="0" cellspacing="0"
								cellpadding="0" style="margin: 0 auto;">
								<tr>
									<td mc:edit="text003" align="center" class="text_color_282828"
										style="color: #282828; font-size: 15px; line-height: 2; font-weight: 500; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text" style="line-height: 2;">
											<span class="text_container">
												<multiline>
													{{.Message}}
												</multiline>
											</span>
										</div>
									</td>
								</tr>
								<!-- horizontal gap -->
								<tr>
									<td height="50"></td>
								</tr>
	
								<tr>
									<td align="center">
										<!-- button -->
										<table class="button_bg_color_303f9f bg_color_303f9f" bgcolor="#303f9f" width="225"
											height="50" align="center" border="0" cellpadding="0" cellspacing="0"
											style="background-color:#303f9f; border-radius:3px;">
											<tr>
												<td mc:edit="text004" align="center" valign="middle"
													style="color: #ffffff; font-size: 16px; font-weight: 600; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;"
													class="text_color_ffffff">
													<div class="editable-text">
														<span class="text_container">
															<multiline>
																<a href="{{.ButtonLink}}"
																	style="text-decoration: none; color: #ffffff;">{{.ButtonTitle}}</a>
															</multiline>
														</span>
													</div>
												</td>
											</tr>
										</table><!-- END button -->
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="25"></td>
								</tr>
							</table><!-- END container_400 -->
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="60"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1" width="600" align="center" border="0" cellspacing="0" cellpadding="0"
					style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="40"></td>
					</tr>
	
					<tr>
						<td>
							<!--  column-1 -->
							<table class="table1-2" width="600" align="left" border="0" cellspacing="0" cellpadding="0">
								<tr>
									<td mc:edit="text006" align="left" class="center_content text_color_929292"
										style="color: #929292; font-size: 14px; line-height: 2; font-weight: 400; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text" style="line-height: 2;">
											<span class="text_container">
												<multiline>
													Recebeu esta mensagem porque está cadastrado no portal <a
														href="https://{{.Domain}}"
														style="color: #303f9f;text-decoration: none;">
														{{.Domain}}</a>
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="20"></td>
								</tr>
	
								<tr>
									<td mc:edit="text008" align="left" class="center_content"
										style="font-size: 14px;font-weight: 400; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text">
											<span class="text_container">
												<multiline>
													<a href="{{.UnsubscribeLink}}" class="text_color_929292"
														style="color:#929292; text-decoration: none; display: block;">Deixar
														de receber mensagens</a>
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- margin-bottom -->
								<tr>
									<td height="30"></td>
								</tr>
							</table><!-- END column-1 -->
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="70"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	</table><!-- END wrapper -->`

	MessageTemplate = `<!-- Section-1 -->
	<table class="table_full editable-bg-color bg_color_e6e6e6 editable-bg-image" bgcolor="#e6e6e6" width="100%"
		align="center" mc:repeatable="castellab" mc:variant="Header" cellspacing="0" cellpadding="0" border="0">
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_303f9f" bgcolor="#303f9f" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="25"></td>
					</tr>
					<tr>
						<td>
							<!-- Inner container -->
							<table class="table1" width="520" align="center" border="0" cellspacing="0" cellpadding="0"
								style="margin: 0 auto;">
								<!-- horizontal gap -->
								<tr>
									<td height="40"></td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="20"></td>
								</tr>
	
								<tr>
									<td mc:edit="text001" align="center" class="text_color_ffffff"
										style="color: #ffffff; font-size: 30px; font-weight: 700; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text">
											<span class="text_container">
												<multiline>
													{{.Subject}}
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="30"></td>
								</tr>
							</table><!-- END inner container -->
						</td>
					</tr>
					<!-- padding-bottom -->
					<tr>
						<td height="60"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_ffffff" bgcolor="#ffffff" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="60"></td>
					</tr>
	
					<tr>
						<td>
							<!-- container_400 -->
							<table class="container_400" align="center" width="400" border="0" cellspacing="0"
								cellpadding="0" style="margin: 0 auto;">
								<tr>
									<td mc:edit="text003" align="center" class="text_color_282828"
										style="color: #282828; font-size: 15px; line-height: 2; font-weight: 500; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text" style="line-height: 2;">
											<span class="text_container">
												<multiline>
													{{.Message}}
												</multiline>
											</span>
										</div>
									</td>
								</tr>
								<!-- horizontal gap -->
								<tr>
									<td height="50"></td>
								</tr>
							</table><!-- END container_400 -->
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="40"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1" width="600" align="center" border="0" cellspacing="0" cellpadding="0"
					style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="40"></td>
					</tr>
	
					<tr>
						<td>
							<!--  column-1 -->
							<table class="table1-2" width="600" align="left" border="0" cellspacing="0" cellpadding="0">
								<tr>
									<td mc:edit="text006" align="left" class="center_content text_color_929292"
										style="color: #929292; font-size: 14px; line-height: 2; font-weight: 400; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text" style="line-height: 2;">
											<span class="text_container">
												<multiline>
													Recebeu esta mensagem porque está cadastrado no portal <a
														href="https://{{.Domain}}"
														style="color: #303f9f;text-decoration: none;">
														{{.Domain}}</a>
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="20"></td>
								</tr>
	
								<tr>
									<td mc:edit="text008" align="left" class="center_content"
										style="font-size: 14px;font-weight: 400; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text">
											<span class="text_container">
												<multiline>
													<a href="{{.UnsubscribeLink}}" class="text_color_929292"
														style="color:#929292; text-decoration: none; display: block;">Deixar
														de receber mensagens</a>
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- margin-bottom -->
								<tr>
									<td height="30"></td>
								</tr>
							</table><!-- END column-1 -->
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="70"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	</table><!-- END wrapper -->`

	ContactTemplate = `<!-- Section-1 -->
	<table class="table_full editable-bg-color bg_color_e6e6e6 editable-bg-image" bgcolor="#e6e6e6" width="100%"
		align="center" mc:repeatable="castellab" mc:variant="Header" cellspacing="0" cellpadding="0" border="0">
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_303f9f" bgcolor="#303f9f" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="25"></td>
					</tr>
					<tr>
						<td>
							<!-- Inner container -->
							<table class="table1" width="520" align="center" border="0" cellspacing="0" cellpadding="0"
								style="margin: 0 auto;">
								<!-- horizontal gap -->
								<tr>
									<td height="40"></td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="20"></td>
								</tr>
	
								<tr>
									<td mc:edit="text001" align="center" class="text_color_ffffff"
										style="color: #ffffff; font-size: 30px; font-weight: 700; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text">
											<span class="text_container">
												<multiline>
													{{.Subject}}<br>
													<span
														style="font-size: 18px; font-weight: 400; line-height: 20px;"><br><a
															style="color: #ffffff;"
															href="mailto:{{.Email}}">{{.Email}}</a></span>
												</multiline>
											</span>
										</div>
									</td>
								</tr>
	
								<!-- horizontal gap -->
								<tr>
									<td height="10"></td>
								</tr>
							</table><!-- END inner container -->
						</td>
					</tr>
					<!-- padding-bottom -->
					<tr>
						<td height="60"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1 editable-bg-color bg_color_ffffff" bgcolor="#ffffff" width="600" align="center"
					border="0" cellspacing="0" cellpadding="0" style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="60"></td>
					</tr>
	
					<tr>
						<td>
							<!-- container_400 -->
							<table class="container_400" align="center" width="400" border="0" cellspacing="0"
								cellpadding="0" style="margin: 0 auto;">
								<tr>
									<td mc:edit="text003" align="center" class="text_color_282828"
										style="color: #282828; font-size: 15px; line-height: 2; font-weight: 500; font-family: lato, Helvetica, sans-serif; mso-line-height-rule: exactly;">
										<div class="editable-text" style="line-height: 2;">
											<span class="text_container">
												<multiline>
													{{.Message}}
												</multiline>
											</span>
										</div>
									</td>
								</tr>
								<!-- horizontal gap -->
								<tr>
									<td height="50"></td>
								</tr>
							</table><!-- END container_400 -->
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="60"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	
		<tr>
			<td>
				<!-- container -->
				<table class="table1" width="600" align="center" border="0" cellspacing="0" cellpadding="0"
					style="margin: 0 auto;">
					<!-- padding-top -->
					<tr>
						<td height="40"></td>
					</tr>
	
					<tr>
						<td>
	
						</td>
					</tr>
	
					<!-- padding-bottom -->
					<tr>
						<td height="70"></td>
					</tr>
				</table><!-- END container -->
			</td>
		</tr>
	</table><!-- END wrapper -->`
)
